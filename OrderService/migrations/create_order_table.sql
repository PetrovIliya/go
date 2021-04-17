CREATE TABLE `orderservice`.`order` (
    `order_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `cost` FLOAT UNSIGNED NOT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` VARCHAR(45) NULL,
    PRIMARY KEY (`order_id`),
    UNIQUE INDEX `order_id_UNIQUE` (`order_id` ASC) VISIBLE
);