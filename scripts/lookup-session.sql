SELECT token FROM sessions WHERE token = '%s' AND expires > DATETIME('now')
