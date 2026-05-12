pipeline {
    agent any

    environment {
        IMAGE_NAME = 'dickaf/shipping-service'
        IMAGE_TAG  = "${BUILD_ID}"
    }

    stages {

        stage('1. Checkout Repo') {
            steps {
                echo 'Checkout source code...'
                checkout scm
            }
        }

        stage('2. Verify Environment') {
            steps {
                bat '''
                echo ===== ENVIRONMENT =====
                go version
                docker version
                kubectl version --client
                '''
            }
        }

        stage('3. Unit Test') {
            steps {
                bat 'go mod tidy'
                bat 'go test -v ./...'
            }
        }

        stage('4. Go Vet') {
            steps {
                bat 'go vet ./...'
            }
        }

        stage('5. Build Docker Image') {
            steps {
                bat "docker build -t ${IMAGE_NAME}:${IMAGE_TAG} ."
            }
        }

        stage('6. Docker Login') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'dockerhub-login',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {

                    bat 'docker login -u %DOCKER_USER% -p %DOCKER_PASS%'
                }
            }
        }

        stage('7. Push Docker Image') {
            steps {
                bat "docker push ${IMAGE_NAME}:${IMAGE_TAG}"
            }
        }

        stage('8. Deploy Kubernetes') {
            steps {
                bat 'kubectl apply -f deployment.yaml'
            }
        }

        stage('9. Verify Deployment') {
            steps {
                bat 'kubectl rollout status deployment/shipping-deployment'
            }
        }
    }

    post {
        always {
            echo 'Pipeline FINISHED'
        }

        success {
            echo 'Pipeline SUCCESS'
        }

        failure {
            echo 'Pipeline FAILED'
        }
    }
}