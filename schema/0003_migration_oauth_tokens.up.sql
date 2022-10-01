CREATE TABLE IF NOT EXISTS oauth_tokens(
	CREATE TABLE oauth_tokens(
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		code VARCHAR(255),
		access VARCHAR(255) NOT NULL,
		refresh VARCHAR(255) NOT NULL,
		data TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		KEY access_k(access),
		KEY refresh_k (refresh),
		KEY expired_at_k (expired_at),
		KEY code_k (code)
	)
);
