import config from './config.json'

const { UserServiceClient } = require('./prpc/UserServiceClientPb.ts')

var serverAddress = config.web_grpc_address

// withCredentials: true for cookie
var userService = new UserServiceClient(serverAddress, null, { withCredentials: true });

export default userService;
export { serverAddress }