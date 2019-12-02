CONTAINER_NAME="yendo"

node {

  stage('Checkout') {
    checkout scm
  }

  try {

    stage('Build: MySQL') {
      sh "docker run --rm --name yendo-mysql -it -e MYSQL_ALLOW_EMPTY_PASSWORD=yes -p 3306:3306 mariadb"
    }

    stage('Test: Go 1.13') {
      sh "docker run --name yendo --rm -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_NAME=mysql golang:1.13 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

    stage('Test: Go 1.12') {
      sh "docker stop yendo"
      sh "docker run --name yendo --rm -v \$(pwd):/go/src/github.com/jamiefdhurst/yendo -e DB_HOST=host.docker.internal -e DB_USER=root -e DB_NAME=mysql golang:1.12 cat"
      sh "docker exec yendo go get github.com/go-sql-driver/mysql"
      sh "docker exec yendo go test github.com/jamiefdhurst/yendo"
    }

  } finally {

    stage('Cleanup') {
      sh "docker stop yendo"
      sh "docker stop yendo-mysql"
      deleteDir()
    }

  }

}
