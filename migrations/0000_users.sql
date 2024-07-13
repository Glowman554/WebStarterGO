CREATE TABLE `users` (
	`username` text PRIMARY KEY NOT NULL,
	`password_hash` text NOT NULL
);

CREATE TABLE `sessions` (
	`token` text PRIMARY KEY NOT NULL,
	`username` text NOT NULL,
	`creation_date` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	FOREIGN KEY (`username`) REFERENCES `users`(`username`) ON UPDATE cascade ON DELETE cascade
);
