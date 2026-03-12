package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	ds "team-task-manager/internal/datastruct"
	"team-task-manager/internal/supports"
)

func getStatusCode(s string) int {
	switch s {
	case ds.StatusNotFound:
		return http.StatusNotFound
	case ds.StatusServiceError:
		return http.StatusInternalServerError
	case ds.StatusAlreadyExists:
		return http.StatusConflict
	case ds.StatusWrongLoginOrPassword:
		return http.StatusUnauthorized
	case ds.StatusInvalidToken:
		return http.StatusUnauthorized
	case ds.StatusSessionReset:
		return http.StatusUnauthorized
	}

	return http.StatusOK
}

func extractJsonBody[ReqT any](r *http.Request, v ReqT) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&v); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func writeJsonResponse[RespT IWithStatus](w *http.ResponseWriter, resp RespT) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&resp); err != nil {
		return err
	}

	(*w).Header().Set(contentLenKey, strconv.Itoa(len(buf.Bytes())))
	(*w).Header().Set(contentTypeKey, appJSONValue)

	code := getStatusCode(resp.GetStatus())
	(*w).WriteHeader(code)
	_, err := (*w).Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func extractSchemaQuery[ReqT IWithStatus](r *http.Request, v ReqT) error {
	if err := schemaDecoder.Decode(&v, r.URL.Query()); err != nil && err != io.EOF {
		return err
	}

	return nil
}

func Exec[ReqT any, RespT IWithStatus](a ExecArgs[ReqT, RespT]) {
	var req ReqT

	if err := a.requestExtractor(a.httpRequest, &req); err != nil {
		msg := "failed extracting request"
		a.api.logger.ErrorKV(msg, "error", err.Error())

		resp := ds.Status{Message: supports.Concat(msg, ": ", err.Error())}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}

		return
	}

	if err := supports.StructValidator().Struct(&req); err != nil {
		msg := "failed validating request"
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)

		resp := ds.Status{Message: supports.Concat(msg, ": ", err.Error())}
		err = writeJsonResponse(a.httpResponse, resp)
		if err != nil {
			a.api.logger.ErrorKV("failed write response",
				"error", err.Error(), "response", resp)
		}
		return
	}

	resp := a.serviceFunc(&req)
	if resp == nil {
		resp := ds.Status{Message: ds.StatusServiceError}
		msg := "failed execute request on service"
		a.api.logger.ErrorKV(msg, "error", "service return no response", "request", req)
		err := writeJsonResponse(a.httpResponse, resp)
		if err != nil {

		}
		return
	}

	if err := a.responseWriter(a.httpResponse, resp); err != nil {
		msg := "failed writing response"
		http.Error(*a.httpResponse, msg, http.StatusInternalServerError)
		a.api.logger.ErrorKV(msg, "error", err.Error(), "request", req)
	}
}
