# 1. Introduction
This project features 2 different component; Data exploration done in a very short python script, and path search optimization with Golang

Data exploration is done with Python because it is very easy to manipulate small datasets, as well as for visualization.

However, for computation heavy tasks such as path search optimization, it is highly recommended to use a fast language that can also handle concurrency. Between the range of static typed languages from low level C to higher level JVM languages (Java, Scala), Golang is somewhere in the middle where it is low enough to be close to the metal, yet not too verbose like JVM languages, and has a small overhead(light binary footprint) and very fast compilation time. Hence, it is used for implementing path search algorithms.

# 2. Setup
`Note: This whole setup was done in Ubuntu 18.04. You may use a similar *nix system`
## 2.1. Python
Make sure you are using Python 3.6X. Open a terminal in current project folder, run the following to install python dependencies.
```sh
$ pip3 install ./requirements.txt
```
## 2.2. Golang
Setting up Golang is trickier. If you do not wish to rebuild the golang project, you can choose to run the binary directly, assuming everythig works smoothly for you.
1st argument is the name of cleaned x,y points. 2nd, 3rd arguments are X,Y coordinates of starting point respectively. 4th argument is an integer depicting the number of paths you want to generate.
```sh
$ ./yojee cleaned.csv 11.552931 104.933636 4
```
If you do not have Golang installed, you may want to use 1) docker method instead to avoid setting up environment variables. However, a Docker container does not have the luxury of accessing an external file easily which means you will need know how to mount the CSV file to the docker container to run your own file.

Either way, both methods involve recompiling the binary hence there may be some complications involved.

### Method 1) Docker
This has been tested with Docker version 18.09.1-rc1, build bca0068
1) Open a terminal in current project folder
2) Build the docker image
```sh
$ sudo docker build -t my_yojee_test .
```

### Method 2) Golang
1) If you have Golang installed and $GOPATH setup properly, you may place this whole project under $GOPATH/github.com/yojee (or under /home/username/go/src/github.com/yojee assuming you are using default $GOPATH)

2) Run the following
```sh
$ go get -d -v ./...
$ go install -v ./...
$ go build .
```
3. At this point, you should see a binary named "yojee" if everything works smoothly.

# 3. Running the project
There are 2 parts to this project; exploration in Python and path search in Golang

## 3.1. Data exploration in python
Full details of the script is in documentation.pdf. To run the script, open a terminal in project folder, run 
```sh
$ python explore.py
```
Once the script finishes, it should created a file named 'cleaned.csv' for cleaned data. Cleaning details are in documentation.pdf. You should also see a Plotly plot in your default browser.

## 3.2. Path search in Golang
### 3.2.1. Binary method
If you can run the binary "yojee" sucessfully, simply open a terminal in project folder and run
```sh
$ ./yojee cleaned.csv 11.552931 104.933636 4
```
1st argument is the name of cleaned x,y points. 2nd, 3rd arguments are X,Y coordinates of starting point respectively. 4th argument is an integer depicting the number of paths you want to generate.

### 3.2.2. Docker Alternative
If you have to use docker, remember this "cleaned.csv" is an existing file in the image itself unless you mount yours. This means for this docker example, you can only run "cleaned.csv".
```sh
sudo docker run --rm my_yojee_test cleaned.csv 11.552931 104.933636 4
```

### 3.2.3. Understanding output
There are 2 parts to the output. The first part is the output from partition optimization, while the second part is output of nearest neighbour path search for each path (Refer to documentation.pdf). Since we used "4" for this example, we should see 4 different paths. 
```sh

3,8,215,134,iter: 0
37,91,103,129,iter: 1
63,93,113,91,iter: 2
70,96,101,93,iter: 3
87,99,91,83,iter: 4
...
86,91,93,90,iter: 97
86,92,88,94,iter: 98
88,90,93,89,iter: 99
Zero centered nearest neighbour path of worker 0
-9.880066e-004 -1.525879e-003
-9.880066e-004 -1.525879e-003
-4.953384e-003 -1.075745e-003
...
Zero centered nearest neighbour path of worker 1
+4.066467e-003 -2.677917e-003
+6.077766e-003 -4.776001e-003
+7.231712e-003 -4.997253e-003
...
Zero centered nearest neighbour path of worker 2
+2.108574e-003 -2.861023e-003
+2.126694e-003 -3.242493e-003
+2.126694e-003 -3.242493e-003
...
Zero centered nearest neighbour path of worker 3
+1.392365e-004 -3.585815e-003
+7.171631e-004 -7.148743e-003
-1.049042e-004 -8.567810e-003
...
```
`Note: This solution is incomplete. Further algorithm designs will require time to find 3rd party implementations and possibly modifying them, or rewriting them grounds up in an optimized manner. Hence, high level designs of these algorithms are discussed in the documentation.pdf instead.`