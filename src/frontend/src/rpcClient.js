import config from './config.json'

const { UserServiceClient } = require('./prpc/UserServiceClientPb.ts')

var serverAddress
try {
  const localConfig = require('./local_config.json')
  serverAddress = localConfig.web_grpc_address
} catch (err) {
  serverAddress = config.web_grpc_address
}

// withCredentials: true for cookie
var userService = new UserServiceClient(serverAddress, null, { withCredentials: true });

export default userService;
export { serverAddress }