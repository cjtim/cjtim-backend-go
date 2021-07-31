"use strict";

// @see https://docs.mongodb.com/manual/tutorial/write-scripts-for-the-mongo-shell/

var MONGODB_URI = "mongodb://mongodb:mongodb@mongodb";

var db = connect(MONGODB_URI);

db = db.getSiblingDB("production");

var collections = db.createCollection("binance");

var binance = db.getCollection("binance");
binance.insertOne({
  usrName: "John Doe",
  usrDept: "Sales",
  usrTitle: "Executive Account Manager",
  authLevel: 4,
  authDept: ["Sales", "Customers"],
});

printjson(collections);
