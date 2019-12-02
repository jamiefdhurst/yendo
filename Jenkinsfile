CONTAINER_NAME="yendo"

node {

  stage('Checkout') {
    checkout scm
  }

  try {

    stage('Build: MySQL') {
      sh "docker run --rm --name yendo-mysql -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -d mariadb"
      sh "sleep 120"
    }

    stage('Test: Go 1.13') {
      sh "docker run --name yendo --rm --network container:yendo-mysql -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.13 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

    stage('Test: Go 1.12') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm --network container:yendo-mysql -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=127.0.0.1 -e DB_USER=root -e DB_NAME=mysql -d -i golang:1.12 cat"
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
