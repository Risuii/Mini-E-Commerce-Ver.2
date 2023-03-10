CREATE TABLE `ecommerce`.`users` (
    `ID` INT NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(255) NULL,
    `password` VARCHAR(255) NULL,
    `email` VARCHAR(255) NULL,
    `address` VARCHAR(255) NULL,
    `created_at` DATE NULL DEFAULT (now()),
    `update_at` DATE NULL DEFAULT (now()),
    PRIMARY KEY (`ID`)
);

CREATE TABLE `ecommerce`.`stores` (
    `ID` INT NOT NULL AUTO_INCREMENT,
    `userID` INT NOT NULL,
    `nameStore` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255),
    `created_at` DATE NULL DEFAULT (now()),
    `update_at` DATE NULL DEFAULT (now()),
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`userID`) REFERENCES users(`ID`)
);

CREATE TABLE `ecommerce`.`items` (
    `ID` INT NOT NULL AUTO_INCREMENT,
    `storeID` INT NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `description` VARCHAR(255) NULL,
    `quantity` INT NULL,
    `created_at` DATE NULL DEFAULT (now()),
    `update_at` DATE NULL DEFAULT (now()),
    PRIMARY KEY (`ID`),
    FOREIGN KEY (`storeID`) REFERENCES stores(`ID`)
);