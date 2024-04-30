package connection_manager

// WebSocket connections are limited to a maximum of 15 connection requests per minute.
// If this limit is exceeded, new connection requests will be temporarily blocked for 60 seconds.
//type BtcturkConnectionManager struct {
//	subscriptions map[string]*websocket.BtcturkWsConnection //
//	subscriptionsLock
//}
