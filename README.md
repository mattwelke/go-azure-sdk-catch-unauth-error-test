# go-azure-sdk-catch-unauth-error

Testing out handling errors from the Go Azure SDK.

Example, catching when the SDK is authenticated but not authorized (via role assignments):

```go
client unauthorized (ensure service principal used to authenticate Azure plugin has permissions for operations "Microsoft.Authorization/roleAssignments/read" and "Microsoft.Authorization/denyAssignments/read" via role assignments and that no deny assignments forbit them): GET https://management.azure.com/subscriptions/<sub-id>/providers/Microsoft.Authorization/roleAssignments
--------------------------------------------------------------------------------
RESPONSE 403: 403 Forbidden
ERROR CODE: AuthorizationFailed
--------------------------------------------------------------------------------
{
  "error": {
    "code": "AuthorizationFailed",
    "message": "The client '<client-id>' with object id '<client-id>' does not have authorization to perform action 'Microsoft.Authorization/roleAssignments/read' over scope '/subscriptions/<sub-id>' or the scope is invalid. If access was recently granted, please refresh your credentials."
  }
}
--------------------------------------------------------------------------------
```
