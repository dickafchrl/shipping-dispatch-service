pipeline {
    agent any

    environment {
        // Ganti dengan username Docker Hub Anda
        DOCKER_IMAGE = "dickaf/shipping-service:latest"
    }

    stages {
        stage('1. Checkout Repo') {
            steps {
                checkout scm
            }
        }

        stage('2. Unit Tests') {
            steps {
                // Menjalankan test tanpa tag functional (hanya logic)
                // Kita bungkus catchError agar pipeline tetap lanjut ke tahap infrastruktur meskipun RED
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'go test -v ./internal/service/... ./internal/handler/http/...'
                }
            }
        }

        stage('3. Lint & Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        stage('4. Build Image (Lokal Jenkins)') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
            }
        }

        stage('5. Functional Tests') {
            steps {
                // 1. Nyalakan DB pendukung
                sh 'docker-compose up -d postgres mongodb'
                
                // 2. Beri jeda agar DB siap menerima koneksi
                echo "Waiting for databases to be ready..."
                sleep 15
                
                // 3. Jalankan functional test dengan Build Tags
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'go test -v -tags=functional ./internal/repository/...'
                }
                
                // 4. Matikan DB setelah test selesai
                sh 'docker-compose down'
            }
        }

        stage('6. Push Image') {
            steps {
                // Membutuhkan credential dengan ID 'docker-hub-id' yang sudah diset di Jenkins
                withCredentials([usernamePassword(credentialsId: 'docker-hub-id', passwordVariable: 'DOCKER_PASS', usernameVariable: 'DOCKER_USER')]) {
                    sh "echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin"
                    sh "docker push ${DOCKER_IMAGE}"
                }
            }
        }

        stage('7. Deploy to Kubernetes') {
            steps {
                // Membutuhkan credential dengan ID 'kubeconfig-id' yang sudah diset di Jenkins
                withKubeConfig([credentialsId: 'kubeconfig-id']) {
                    sh 'kubectl apply -f k8s/deployment.yaml'
                    sh 'kubectl apply -f k8s/service.yaml'
                }
            }
        }

        stage('8. Verify Rollout') {
            steps {
                sh 'kubectl rollout status deployment/shipping-service-deployment'
            }
        }
    }

    post {
        always {
            echo "Pipeline finished. Check test results in the stages above."
        }
    }
}