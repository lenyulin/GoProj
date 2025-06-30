package ORM

//ent 是 Facebook 开源的一款 Go 语言实体框架，是一款简单而强大的用于建模和查询数据的 ORM 框架。
//link: https://www.liwenzhou.com/posts/Go/ent/

//Bun 是一个 SQL 优先的 Golang ORM（对象关系映射），支持 PostgreSQL、MySQL、MSSQL和SQLite。它旨在提供一种简单高效的数据库使用方法，同时利用 Go 的类型安全性并减少重复代码。

//与 GORM 比较
//相比于使用 GORM 框架，使用 Bun 可以更容易地编写复杂查询。Bun 可以更好地集成特定数据库的功能，例如 PostgreSQL arrays。通常 Bun 的速度也更快。
//
//Bun 不支持自动迁移、优化器/索引/注释提示和数据库解析器等 GORM 常用功能。
//
//与 ent 比较
//使用 Bun 的时候你可以利用以前使用 SQL DBMS 和 Go 的经验来编写快速、习惯化的代码。Bun 的目标是帮助你编写 SQL 而不是取代或隐藏它。
//使用 ent 框架则无法利用之前的经验，因为 ent 提供了一种全新/不同的方法来编写 Go 应用程序，你只能遵循它的规则。
