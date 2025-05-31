package main

func Handle(req Request, store *Store) (string, error) {
	switch req.Command {
	case CMD_GET:
		key := req.Args[0]
		value, ok := store.Get(key)
		if !ok {
			return "", &RequestError{Message: "key not found !"}
		}
		return value, nil
	case CMD_SET:
		key, value := req.Args[0], req.Args[1]
		store.Set(key, value)
		return "OK", nil
	case CMD_DEL:
		key := req.Args[0]
		store.Del(key)
		return "OK", nil
	case CMD_HELP:
		helpMessage := "Supported commands:\n" +
			"GET <key>\n" +
			"SET <key> <value>\n" +
			"DEL <key>\n" +
			"HELP\n" +
			"QUIT"
		return helpMessage, nil
	case CMD_QUIT:
		return "Goodbye!", nil
	case CMD_UNKNOWN:
		return "", &RequestError{Message: "Unknown Command!"}
	default:
		return "", &RequestError{Message: "Unhandled command type!"}
	}
}
