# sample ad-tracker 
-tracks all kinds of parmaetrs in digital-campaigns
-written in GoLang using Go-Micro Framework
-each feature is written in micro-services and can individually scale
-uses live database (redis), logging database(mongo), seacrh database(Elastic Search), Fast Database(AeroSpike), Data Warehouse(big query)
-API contains urls registerd with service discovery system(Consul)
-Listners are responsible for carrying out business logic 
-Services are used for complex operations registered in Consul
-Handlers carry out config operations interacting with Logging DB(mongo) and live data(redis)
-Subscribers are responsible for report generations, cleaning up db, bulk processing on queue data,upload to big query
and generating match for repeated users.

