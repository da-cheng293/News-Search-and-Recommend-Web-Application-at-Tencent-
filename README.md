That is a news search and recommend web application. First, all news will be collected from the news provider through their api and web crawler techniques.  after processing those data to the same protocol, they will be pushed to kafka waiting consumer to consume. I will also design some filters at the consumer side such as sensitive news check, removing duplication using md5 implemented in redis. Then, it will be  classfied and labeled through spark. all the processed data will be store in the elasticsearch database and waiting the web server to use for searching purpose. all the above development will be depolyed using docker techniques. 
