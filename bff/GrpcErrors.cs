using Grpc.Core;

namespace Bff;

public static class GrpcErrors
{
    // gRPC has its own status codes. The browser speaks HTTP, so the boundary translates them.
    public static IResult ToProblem(this RpcException ex) => Results.Problem(
        title: ex.StatusCode.ToString(),
        detail: ex.Status.Detail,
        statusCode: ex.StatusCode switch
        {
            StatusCode.InvalidArgument => StatusCodes.Status400BadRequest,
            StatusCode.NotFound => StatusCodes.Status404NotFound,
            StatusCode.FailedPrecondition => StatusCodes.Status409Conflict,
            StatusCode.DeadlineExceeded => StatusCodes.Status504GatewayTimeout,
            StatusCode.Unavailable => StatusCodes.Status503ServiceUnavailable,
            _ => StatusCodes.Status502BadGateway,
        });
}
