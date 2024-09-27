SELECT token, expires, DATETIME('now') FROM sessions WHERE token = $1 AND expires > DATETIME('now')
