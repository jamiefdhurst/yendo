CONTAINER_NAME="yendo"

node {

  stage('Checkout') {
    checkout scm
  }

  try {

    stage('Build: MySQL') {
      sh "docker run --rm --name yendo-mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mariadb"
      sh "docker run --rm --network container:yendo-mysql willwill/wait-for-it 127.0.0.1:3306 -t 600"
      sh "sleep 15"
    }

    stage('Test: Go 1.17') {
      sh "docker run --name yendo --rm --network container:yendo-mysql -w /go/src/github.com/jamiefdhurst/yendo -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.17 cat"
      sh "docker exec yendo go mod download"
      sh "docker exec yendo go test ."
    }

    stage('Test: Go 1.16') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm --network container:yendo-mysql -w /go/src/github.com/jamiefdhurst/yendo -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.16 cat"
      sh "docker exec yendo go mod download"
      sh "docker exec yendo go test ."
    }

    stage('Test: Go 1.15') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm --network container:yendo-mysql -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.15 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

    stage('Test: Go 1.14') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm --network container:yendo-mysql -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.14 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

    stage('Test: Go 1.13') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm --network container:yendo-mysql -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.13 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

  } finally {

    stage('Cleanup') {
      sh "docker stop yendo-mysql"
      sh "docker stop yendo"
      deleteDir()
    }

  }

}
