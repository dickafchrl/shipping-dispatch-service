pipeline {
    agent any

    environment {
        DOCKER_USER = "dickaf" 
        IMAGE_NAME = "shipping-service"
        DOCKER_IMAGE = "${DOCKER_USER}/${IMAGE_NAME}:latest"
    }

    stages {
        stage('1. Setup Workspace') {
            steps {
                // Pastikan folder bin ada dan docker-compose terinstall permanen di container jenkins
                sh '''
                    if ! command -v docker-compose &> /dev/null; then
                        curl -L "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
                        chmod +x /usr/local/bin/docker-compose
                    fi
                '''
            }
        }

        stage('2. Unit Tests') {
            steps {
                // Kita tidak memanggil 'go' di luar, semua di dalam kontainer.
                // Kita gunakan './...' yang secara otomatis akan mengabaikan file dengan tag 'functional' 
                // SELAMA kita tidak menambahkan '-tags=functional'
                sh '''
                    docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
                    -v $(pwd):/app -w /app golang:alpine \
                    sh -c "go test -v ./internal/service/... ./internal/handler/http/..."
                '''
            }
        }

        stage('3. Build Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
            }
        }

        stage('4. Functional Tests') {
            steps {
                // Bersihkan container lama yang namanya bentrok
                sh 'docker-compose down || true'
                sh 'docker rm -f shipping_mongo shipping_postgres || true'
                
                // Jalankan DB baru
                sh 'docker-compose up -d postgres mongodb'
                echo "Waiting for DBs..."
                sleep 20
                
                // Jalankan functional test di dalam container Go
                sh 'docker run --rm --network host -v $(pwd):/app -w /app golang:alpine go test -v -tags=functional ./internal/repository/...'
                
                sh 'docker-compose down'
            }
        }

        stage('5. Push to Docker Hub') {
            steps {
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