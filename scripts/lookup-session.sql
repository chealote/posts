SELECT token FROM sessions WHERE token = $1 AND expires > DATETIME('now')
