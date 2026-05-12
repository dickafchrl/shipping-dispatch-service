pipeline {
    agent any

    environment {
        // SESUAIKAN dengan username Docker Hub yang baru saja Anda buat
        DOCKER_USER = "dickaf" 
        IMAGE_NAME = "shipping-service"
        DOCKER_IMAGE = "${DOCKER_USER}/${IMAGE_NAME}:latest"
    }

    stages {
        stage('1. Setup Workspace') {
            steps {
                // Memastikan docker-compose tersedia di dalam container Jenkins
                sh 'docker-compose --version || (curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose)'
            }
        }

        stage('2. Unit Tests') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    // Kita asumsikan image golang sudah terpasang atau gunakan container agent
                    sh 'go test -v ./internal/service/... ./internal/handler/http/...'
                }
            }
        }

        stage('3. Build Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
            }
        }

        stage('4. Functional Tests') {
            steps {
                sh 'docker-compose up -d postgres mongodb'
                echo "Waiting for DBs..."
                sleep 20
                
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'go test -v -tags=functional ./internal/repository/...'
                }
                
                sh 'docker-compose down'
            }
        }

        stage('5. Push to Docker Hub') {
            steps {
                // Pakai ID credential 'docker-hub-id' yang Anda buat di Jenkins
                withCredentials([usernamePassword(credentialsId: 'docker-hub-id', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER_ENV')]) {
                    sh "echo \$DOCKER_PASS | docker login -u \$DOCKER_USER_ENV --password-stdin"
                    sh "docker push ${DOCKER_IMAGE}"
                }
            }
        }

        stage('6. Deploy to K8s') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-id']) {
                    sh 'kubectl apply -f k8s/deployment.yaml'
                    sh 'kubectl apply -f k8s/service.yaml'
                }
            }
        }
    }
}