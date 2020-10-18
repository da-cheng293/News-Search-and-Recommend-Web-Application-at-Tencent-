from pyspark import SparkConf, SparkContext 

def word_count():
	conf = SparkConf().setAppName('test').setMaster('spark://localhost:7077') 
	sc = SparkContext(conf=conf) 


	text_file = sc.textFile("news.txt")

	# word count
	counts = text_file.flatMap( lambda line: line.lower().split(" ") ) \
	            .map( lambda word: (word, 1) ) \
	            .reduceByKey( lambda a, b: a + b ) \
	            .sortBy( lambda x: x[1], False )
	output = counts.collect()


	for (word, count) in output:
	    print( "%s: %i" % (word, count) )
	    
	# Stopping Spark Context
	sc.stop()