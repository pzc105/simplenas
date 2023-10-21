import * as config from './config.js'

const { UserServiceClient } = require('./prpc/UserServiceClientPb.ts')

var serverAddress = config.rpc_server

// withCredentials: true for cookie
var userService = new UserServiceClient(serverAddress, null, { withCredentials: true });

export default userService;
export { serverAddress }