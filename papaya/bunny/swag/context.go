package swag

import (
  "github.com/gofiber/fiber/v2"
  "net/http"
  "skfw/papaya/koala/kornet"
)

type SwagContext struct {
  *fiber.Ctx
  event    bool
  data     any
  body     []byte
  kReq     *kornet.Request
  kRes     *kornet.Response
  finished bool
}

type SwagContextImpl interface {

  // create event if run on a task

  Event() bool    // if a task running, event is true
  Target() any    // Event
  Solve(data any) // --> Target()

  // replace data body

  Modify(body []byte)

  // wrapping from binding

  Body() []byte
  Status(status int) *SwagContext

  // stop a task and send quickly

  Prevent()
  Revoke() bool

  // Kornet binding for any purpose

  Kornet() (*kornet.Request, *kornet.Response)
  Bind(req *kornet.Request, res *kornet.Response)

  // wrapping http status + message + data

  Continue(result *kornet.Result) error
  SwitchingProtocols(result *kornet.Result) error
  Processing(result *kornet.Result) error
  EarlyHints(result *kornet.Result) error
  OK(result *kornet.Result) error
  Created(result *kornet.Result) error
  Accepted(result *kornet.Result) error
  NonAuthoritativeInfo(result *kornet.Result) error
  NoContent(result *kornet.Result) error
  ResetContent(result *kornet.Result) error
  PartialContent(result *kornet.Result) error
  MultiStatus(result *kornet.Result) error
  AlreadyReported(result *kornet.Result) error
  IMUsed(result *kornet.Result) error
  MultipleChoices(result *kornet.Result) error
  MovedPermanently(result *kornet.Result) error
  Found(result *kornet.Result) error
  SeeOther(result *kornet.Result) error
  NotModified(result *kornet.Result) error
  UseProxy(result *kornet.Result) error
  TemporaryRedirect(result *kornet.Result) error
  PermanentRedirect(result *kornet.Result) error
  BadRequest(result *kornet.Result) error
  Unauthorized(result *kornet.Result) error
  PaymentRequired(result *kornet.Result) error
  Forbidden(result *kornet.Result) error
  NotFound(result *kornet.Result) error
  MethodNotAllowed(result *kornet.Result) error
  NotAcceptable(result *kornet.Result) error
  ProxyAuthRequired(result *kornet.Result) error
  RequestTimeout(result *kornet.Result) error
  Conflict(result *kornet.Result) error
  Gone(result *kornet.Result) error
  LengthRequired(result *kornet.Result) error
  PreconditionFailed(result *kornet.Result) error
  RequestEntityTooLarge(result *kornet.Result) error
  RequestURITooLong(result *kornet.Result) error
  UnsupportedMediaType(result *kornet.Result) error
  RequestedRangeNotSatisfiable(result *kornet.Result) error
  ExpectationFailed(result *kornet.Result) error
  MisdirectedRequest(result *kornet.Result) error
  UnprocessableEntity(result *kornet.Result) error
  Locked(result *kornet.Result) error
  FailedDependency(result *kornet.Result) error
  TooEarly(result *kornet.Result) error
  UpgradeRequired(result *kornet.Result) error
  PreconditionRequired(result *kornet.Result) error
  TooManyRequests(result *kornet.Result) error
  RequestHeaderFieldsTooLarge(result *kornet.Result) error
  UnavailableForLegalReasons(result *kornet.Result) error
  InternalServerError(result *kornet.Result) error
  NotImplemented(result *kornet.Result) error
  BadGateway(result *kornet.Result) error
  ServiceUnavailable(result *kornet.Result) error
  GatewayTimeout(result *kornet.Result) error
  HTTPVersionNotSupported(result *kornet.Result) error
  VariantAlsoNegotiates(result *kornet.Result) error
  InsufficientStorage(result *kornet.Result) error
  LoopDetected(result *kornet.Result) error
  NotExtended(result *kornet.Result) error
  NetworkAuthenticationRequired(result *kornet.Result) error

  // wrap a message only

  Message(message string) error

  // empty or nothing

  Empty() error
}

func MakeSwagContext(ctx *fiber.Ctx, event bool) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: event,
    data:  nil,
  }
}

func MakeSwagContextEvent(ctx *fiber.Ctx, data any) *SwagContext {

  return &SwagContext{
    Ctx:   ctx,
    event: true,
    body:  nil,
    data:  data,
  }
}

func (c *SwagContext) Event() bool {

  return c.event
}

func (c *SwagContext) Target() any {

  return c.data
}

func (c *SwagContext) Solve(data any) {

  c.event = true // solve value from key expectation
  c.data = data  // value from expectation
}

func (c *SwagContext) Modify(body []byte) {

  c.body = body
}

// wrapping Body with Modify Body

func (c *SwagContext) Body() []byte {

  if c.body != nil {

    return c.body
  }

  return c.Ctx.Body()
}

// wrapping Status with New SwagContext

func (c *SwagContext) Status(status int) *SwagContext {

  return &SwagContext{
    Ctx:      c.Ctx.Status(status),
    event:    c.event,
    data:     c.data,
    body:     c.body,
    finished: c.finished,
  }
}

func (c *SwagContext) Prevent() {

  c.finished = true
}

func (c *SwagContext) Revoke() bool {

  return c.finished
}

func (c *SwagContext) Kornet() (*kornet.Request, *kornet.Response) {

  return c.kReq, c.kRes
}

func (c *SwagContext) Bind(req *kornet.Request, res *kornet.Response) {

  c.kReq = req
  c.kRes = res

  if req.Body.Size() > 0 {

    c.body = req.Body.ReadAll()
    req.Body.Seek(0)
  }
}

func (c *SwagContext) Continue(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusContinue)
  result.Error = kornet.HttpCheckErrStat(http.StatusContinue)
  return c.Ctx.Status(http.StatusContinue).JSON(result)
}

func (c *SwagContext) SwitchingProtocols(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusSwitchingProtocols)
  result.Error = kornet.HttpCheckErrStat(http.StatusSwitchingProtocols)
  return c.Ctx.Status(http.StatusSwitchingProtocols).JSON(result)
}

func (c *SwagContext) Processing(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusProcessing)
  result.Error = kornet.HttpCheckErrStat(http.StatusProcessing)
  return c.Ctx.Status(http.StatusProcessing).JSON(result)
}

func (c *SwagContext) EarlyHints(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusEarlyHints)
  result.Error = kornet.HttpCheckErrStat(http.StatusEarlyHints)
  return c.Ctx.Status(http.StatusEarlyHints).JSON(result)
}

func (c *SwagContext) OK(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusOK)
  result.Error = kornet.HttpCheckErrStat(http.StatusOK)
  return c.Ctx.Status(http.StatusOK).JSON(result)
}

func (c *SwagContext) Created(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusCreated)
  result.Error = kornet.HttpCheckErrStat(http.StatusCreated)
  return c.Ctx.Status(http.StatusCreated).JSON(result)
}

func (c *SwagContext) Accepted(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusAccepted)
  result.Error = kornet.HttpCheckErrStat(http.StatusAccepted)
  return c.Ctx.Status(http.StatusAccepted).JSON(result)
}

func (c *SwagContext) NonAuthoritativeInfo(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNonAuthoritativeInfo)
  result.Error = kornet.HttpCheckErrStat(http.StatusNonAuthoritativeInfo)
  return c.Ctx.Status(http.StatusNonAuthoritativeInfo).JSON(result)
}

func (c *SwagContext) NoContent(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNoContent)
  result.Error = kornet.HttpCheckErrStat(http.StatusNoContent)
  return c.Ctx.Status(http.StatusNoContent).JSON(result)
}

func (c *SwagContext) ResetContent(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusResetContent)
  result.Error = kornet.HttpCheckErrStat(http.StatusResetContent)
  return c.Ctx.Status(http.StatusResetContent).JSON(result)
}

func (c *SwagContext) PartialContent(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusPartialContent)
  result.Error = kornet.HttpCheckErrStat(http.StatusPartialContent)
  return c.Ctx.Status(http.StatusPartialContent).JSON(result)
}

func (c *SwagContext) MultiStatus(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusMultiStatus)
  result.Error = kornet.HttpCheckErrStat(http.StatusMultiStatus)
  return c.Ctx.Status(http.StatusMultiStatus).JSON(result)
}

func (c *SwagContext) AlreadyReported(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusAlreadyReported)
  result.Error = kornet.HttpCheckErrStat(http.StatusAlreadyReported)
  return c.Ctx.Status(http.StatusAlreadyReported).JSON(result)
}

func (c *SwagContext) IMUsed(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusIMUsed)
  result.Error = kornet.HttpCheckErrStat(http.StatusIMUsed)
  return c.Ctx.Status(http.StatusIMUsed).JSON(result)
}

func (c *SwagContext) MultipleChoices(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusMultipleChoices)
  result.Error = kornet.HttpCheckErrStat(http.StatusMultipleChoices)
  return c.Ctx.Status(http.StatusMultipleChoices).JSON(result)
}

func (c *SwagContext) MovedPermanently(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusMovedPermanently)
  result.Error = kornet.HttpCheckErrStat(http.StatusMovedPermanently)
  return c.Ctx.Status(http.StatusMovedPermanently).JSON(result)
}

func (c *SwagContext) Found(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusFound)
  result.Error = kornet.HttpCheckErrStat(http.StatusFound)
  return c.Ctx.Status(http.StatusFound).JSON(result)
}

func (c *SwagContext) SeeOther(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusSeeOther)
  result.Error = kornet.HttpCheckErrStat(http.StatusSeeOther)
  return c.Ctx.Status(http.StatusSeeOther).JSON(result)
}

func (c *SwagContext) NotModified(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNotModified)
  result.Error = kornet.HttpCheckErrStat(http.StatusNotModified)
  return c.Ctx.Status(http.StatusNotModified).JSON(result)
}

func (c *SwagContext) UseProxy(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUseProxy)
  result.Error = kornet.HttpCheckErrStat(http.StatusUseProxy)
  return c.Ctx.Status(http.StatusUseProxy).JSON(result)
}

func (c *SwagContext) TemporaryRedirect(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusTemporaryRedirect)
  result.Error = kornet.HttpCheckErrStat(http.StatusTemporaryRedirect)
  return c.Ctx.Status(http.StatusTemporaryRedirect).JSON(result)
}

func (c *SwagContext) PermanentRedirect(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusPermanentRedirect)
  result.Error = kornet.HttpCheckErrStat(http.StatusPermanentRedirect)
  return c.Ctx.Status(http.StatusPermanentRedirect).JSON(result)
}

func (c *SwagContext) BadRequest(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusBadRequest)
  result.Error = kornet.HttpCheckErrStat(http.StatusBadRequest)
  return c.Ctx.Status(http.StatusBadRequest).JSON(result)
}

func (c *SwagContext) Unauthorized(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUnauthorized)
  result.Error = kornet.HttpCheckErrStat(http.StatusUnauthorized)
  return c.Ctx.Status(http.StatusUnauthorized).JSON(result)
}

func (c *SwagContext) PaymentRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusPaymentRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusPaymentRequired)
  return c.Ctx.Status(http.StatusPaymentRequired).JSON(result)
}

func (c *SwagContext) Forbidden(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusForbidden)
  result.Error = kornet.HttpCheckErrStat(http.StatusForbidden)
  return c.Ctx.Status(http.StatusForbidden).JSON(result)
}

func (c *SwagContext) NotFound(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNotFound)
  result.Error = kornet.HttpCheckErrStat(http.StatusNotFound)
  return c.Ctx.Status(http.StatusNotFound).JSON(result)
}

func (c *SwagContext) MethodNotAllowed(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusMethodNotAllowed)
  result.Error = kornet.HttpCheckErrStat(http.StatusMethodNotAllowed)
  return c.Ctx.Status(http.StatusMethodNotAllowed).JSON(result)
}

func (c *SwagContext) NotAcceptable(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNotAcceptable)
  result.Error = kornet.HttpCheckErrStat(http.StatusNotAcceptable)
  return c.Ctx.Status(http.StatusNotAcceptable).JSON(result)
}

func (c *SwagContext) ProxyAuthRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusProxyAuthRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusProxyAuthRequired)
  return c.Ctx.Status(http.StatusProxyAuthRequired).JSON(result)
}

func (c *SwagContext) RequestTimeout(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusRequestTimeout)
  result.Error = kornet.HttpCheckErrStat(http.StatusRequestTimeout)
  return c.Ctx.Status(http.StatusRequestTimeout).JSON(result)
}

func (c *SwagContext) Conflict(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusConflict)
  result.Error = kornet.HttpCheckErrStat(http.StatusConflict)
  return c.Ctx.Status(http.StatusConflict).JSON(result)
}

func (c *SwagContext) Gone(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusGone)
  result.Error = kornet.HttpCheckErrStat(http.StatusGone)
  return c.Ctx.Status(http.StatusGone).JSON(result)
}

func (c *SwagContext) LengthRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusLengthRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusLengthRequired)
  return c.Ctx.Status(http.StatusLengthRequired).JSON(result)
}

func (c *SwagContext) PreconditionFailed(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusPreconditionFailed)
  result.Error = kornet.HttpCheckErrStat(http.StatusPreconditionFailed)
  return c.Ctx.Status(http.StatusPreconditionFailed).JSON(result)
}

func (c *SwagContext) RequestEntityTooLarge(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusRequestEntityTooLarge)
  result.Error = kornet.HttpCheckErrStat(http.StatusRequestEntityTooLarge)
  return c.Ctx.Status(http.StatusRequestEntityTooLarge).JSON(result)
}

func (c *SwagContext) RequestURITooLong(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusRequestURITooLong)
  result.Error = kornet.HttpCheckErrStat(http.StatusRequestURITooLong)
  return c.Ctx.Status(http.StatusRequestURITooLong).JSON(result)
}

func (c *SwagContext) UnsupportedMediaType(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUnsupportedMediaType)
  result.Error = kornet.HttpCheckErrStat(http.StatusUnsupportedMediaType)
  return c.Ctx.Status(http.StatusUnsupportedMediaType).JSON(result)
}

func (c *SwagContext) RequestedRangeNotSatisfiable(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusRequestedRangeNotSatisfiable)
  result.Error = kornet.HttpCheckErrStat(http.StatusRequestedRangeNotSatisfiable)
  return c.Ctx.Status(http.StatusRequestedRangeNotSatisfiable).JSON(result)
}

func (c *SwagContext) ExpectationFailed(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusExpectationFailed)
  result.Error = kornet.HttpCheckErrStat(http.StatusExpectationFailed)
  return c.Ctx.Status(http.StatusExpectationFailed).JSON(result)
}

func (c *SwagContext) MisdirectedRequest(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusMisdirectedRequest)
  result.Error = kornet.HttpCheckErrStat(http.StatusMisdirectedRequest)
  return c.Ctx.Status(http.StatusMisdirectedRequest).JSON(result)
}

func (c *SwagContext) UnprocessableEntity(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUnprocessableEntity)
  result.Error = kornet.HttpCheckErrStat(http.StatusUnprocessableEntity)
  return c.Ctx.Status(http.StatusUnprocessableEntity).JSON(result)
}

func (c *SwagContext) Locked(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusLocked)
  result.Error = kornet.HttpCheckErrStat(http.StatusLocked)
  return c.Ctx.Status(http.StatusLocked).JSON(result)
}

func (c *SwagContext) FailedDependency(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusFailedDependency)
  result.Error = kornet.HttpCheckErrStat(http.StatusFailedDependency)
  return c.Ctx.Status(http.StatusFailedDependency).JSON(result)
}

func (c *SwagContext) TooEarly(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusTooEarly)
  result.Error = kornet.HttpCheckErrStat(http.StatusTooEarly)
  return c.Ctx.Status(http.StatusTooEarly).JSON(result)
}

func (c *SwagContext) UpgradeRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUpgradeRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusUpgradeRequired)
  return c.Ctx.Status(http.StatusUpgradeRequired).JSON(result)
}

func (c *SwagContext) PreconditionRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusPreconditionRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusPreconditionRequired)
  return c.Ctx.Status(http.StatusPreconditionRequired).JSON(result)
}

func (c *SwagContext) TooManyRequests(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusTooManyRequests)
  result.Error = kornet.HttpCheckErrStat(http.StatusTooManyRequests)
  return c.Ctx.Status(http.StatusTooManyRequests).JSON(result)
}

func (c *SwagContext) RequestHeaderFieldsTooLarge(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusRequestHeaderFieldsTooLarge)
  result.Error = kornet.HttpCheckErrStat(http.StatusRequestHeaderFieldsTooLarge)
  return c.Ctx.Status(http.StatusRequestHeaderFieldsTooLarge).JSON(result)
}

func (c *SwagContext) UnavailableForLegalReasons(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusUnavailableForLegalReasons)
  result.Error = kornet.HttpCheckErrStat(http.StatusUnavailableForLegalReasons)
  return c.Ctx.Status(http.StatusUnavailableForLegalReasons).JSON(result)
}

func (c *SwagContext) InternalServerError(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusInternalServerError)
  result.Error = kornet.HttpCheckErrStat(http.StatusInternalServerError)
  return c.Ctx.Status(http.StatusInternalServerError).JSON(result)
}

func (c *SwagContext) NotImplemented(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNotImplemented)
  result.Error = kornet.HttpCheckErrStat(http.StatusNotImplemented)
  return c.Ctx.Status(http.StatusNotImplemented).JSON(result)
}

func (c *SwagContext) BadGateway(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusBadGateway)
  result.Error = kornet.HttpCheckErrStat(http.StatusBadGateway)
  return c.Ctx.Status(http.StatusBadGateway).JSON(result)
}

func (c *SwagContext) ServiceUnavailable(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusServiceUnavailable)
  result.Error = kornet.HttpCheckErrStat(http.StatusServiceUnavailable)
  return c.Ctx.Status(http.StatusServiceUnavailable).JSON(result)
}

func (c *SwagContext) GatewayTimeout(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusGatewayTimeout)
  result.Error = kornet.HttpCheckErrStat(http.StatusGatewayTimeout)
  return c.Ctx.Status(http.StatusGatewayTimeout).JSON(result)
}

func (c *SwagContext) HTTPVersionNotSupported(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusHTTPVersionNotSupported)
  result.Error = kornet.HttpCheckErrStat(http.StatusHTTPVersionNotSupported)
  return c.Ctx.Status(http.StatusHTTPVersionNotSupported).JSON(result)
}

func (c *SwagContext) VariantAlsoNegotiates(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusVariantAlsoNegotiates)
  result.Error = kornet.HttpCheckErrStat(http.StatusVariantAlsoNegotiates)
  return c.Ctx.Status(http.StatusVariantAlsoNegotiates).JSON(result)
}

func (c *SwagContext) InsufficientStorage(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusInsufficientStorage)
  result.Error = kornet.HttpCheckErrStat(http.StatusInsufficientStorage)
  return c.Ctx.Status(http.StatusInsufficientStorage).JSON(result)
}

func (c *SwagContext) LoopDetected(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusLoopDetected)
  result.Error = kornet.HttpCheckErrStat(http.StatusLoopDetected)
  return c.Ctx.Status(http.StatusLoopDetected).JSON(result)
}

func (c *SwagContext) NotExtended(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNotExtended)
  result.Error = kornet.HttpCheckErrStat(http.StatusNotExtended)
  return c.Ctx.Status(http.StatusNotExtended).JSON(result)
}

func (c *SwagContext) NetworkAuthenticationRequired(result *kornet.Result) error {

  result.Status = http.StatusText(http.StatusNetworkAuthenticationRequired)
  result.Error = kornet.HttpCheckErrStat(http.StatusNetworkAuthenticationRequired)
  return c.Ctx.Status(http.StatusNetworkAuthenticationRequired).JSON(result)
}

func (c *SwagContext) Message(message string) error {

  return c.OK(kornet.Msg(message, false))
}

func (c *SwagContext) Empty() error {

  return c.Ctx.Status(http.StatusNoContent).Send([]byte{})
}
