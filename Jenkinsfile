pipeline {
    agent any

    environment {
        DOCKER_USER  = "dickaf"
        IMAGE_NAME   = "shipping-service"
        DOCKER_IMAGE = "${DOCKER_USER}/${IMAGE_NAME}:latest"
    }

    stages {

        stage('1. Verify Environment') {
            steps {
                bat '''
                    echo ===== ENVIRONMENT CHECK =====

                    echo Workspace:
                    cd

                    echo.
                    echo ===== GO VERSION =====
                    go version

                    echo.
                    echo ===== DOCKER VERSION =====
                    docker version

                    echo.
                    echo ===== DOCKER COMPOSE VERSION =====
                    docker compose version

                    echo.
                    echo ===== KUBECTL VERSION =====
                    kubectl version --client
                '''
            }
        }

        stage('2. Unit Tests') {
            steps {
                bat '''
                    echo ===== UNIT TEST =====

                    go mod tidy

                    go test -v ./...
                '''
            }
        }

        stage('3. Build Docker Image') {
            steps {
                bat '''
                    echo ===== BUILD IMAGE =====

                    docker build -t %DOCKER_IMAGE% .
                '''
            }
        }

        stage('4. Functional Tests') {
            steps {

                bat '''
                    echo ===== CLEAN OLD CONTAINERS =====

                    docker compose down

                    docker rm -f shipping_mongo
                    docker rm -f shipping_postgres
                '''

                bat '''
                    echo ===== START DATABASE =====

                    docker compose up -d postgres mongodb

                    timeout /t 20
                '''

                bat '''
                    echo ===== FUNCTIONAL TEST =====

                    go test -v -tags=functional ./internal/repository/...
                '''
            }

            post {
                always {
                    bat '''
                        docker compose down
                    '''
                }
            }
        }

        stage('5. Push Docker Image') {
            steps {

                withCredentials([
                    usernamePassword(
                        credentialsId: 'docker-hub-id',
                        usernameVariable: 'DOCKER_USER_ENV',
                        passwordVariable: 'DOCKER_PASS'
                    )
                ]) {

                    bat '''
                        echo ===== DOCKER LOGIN =====

                        echo %DOCKER_PASS% | docker login -u %DOCKER_USER_ENV% --password-stdin

                        docker push %DOCKER_IMAGE%
                    '''
                }
            }
        }

        stage('6. Deploy Kubernetes') {
            steps {

                withKubeConfig([credentialsId: 'kubeconfig-id']) {

                    bat '''
                        echo ===== DEPLOY K8S =====

                        kubectl apply -f k8s\\deployment.yaml

                        kubectl apply -f k8s\\service.yaml

                        kubectl get pods
                    '''
                }
            }
        }
    }

    post {

        success {
            echo 'Pipeline SUCCESS'
        }

        failure {
            echo 'Pipeline FAILED'
        }

        always {
            echo 'Pipeline FINISHED'
        }
    }
}