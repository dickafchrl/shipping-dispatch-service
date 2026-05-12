pipeline {
    agent any

    environment {
        DOCKER_USER  = "dickaf"
        IMAGE_NAME   = "shipping-service"
        DOCKER_IMAGE = "${DOCKER_USER}/${IMAGE_NAME}:latest"
    }

    stages {

        stage('1. Setup Workspace') {
            steps {
                sh '''
                    set -e

                    echo "===== WORKSPACE INFO ====="
                    echo "WORKSPACE: $WORKSPACE"

                    pwd
                    ls -la

                    echo "===== INSTALL DOCKER COMPOSE ====="

                    if ! command -v docker-compose > /dev/null 2>&1; then
                        curl -L \
                        "https://github.com/docker/compose/releases/download/v2.20.2/docker-compose-$(uname -s)-$(uname -m)" \
                        -o /usr/local/bin/docker-compose

                        chmod +x /usr/local/bin/docker-compose
                    fi

                    docker-compose version
                '''
            }
        }

        stage('2. Unit Tests') {
            steps {
                sh '''
                    set -e

                    echo "===== UNIT TEST STAGE ====="

                    echo "HOST WORKSPACE:"
                    pwd
                    ls -la
                    find . -name go.mod

                    docker run --rm \
                        -v $WORKSPACE:/app \
                        -w /app \
                        golang:1.25 \
                        sh -c "
                            set -e

                            echo '===== CONTAINER WORKSPACE ====='

                            pwd
                            ls -la
                            find . -name go.mod

                            if [ ! -f go.mod ]; then
                                echo 'ERROR: go.mod not found'
                                exit 1
                            fi

                            go version

                            go mod tidy

                            go test -v ./...
                        "
                '''
            }
        }

        stage('3. Build Image') {
            steps {
                sh '''
                    set -e

                    echo "===== BUILD IMAGE ====="

                    docker build -t ${DOCKER_IMAGE} .
                '''
            }
        }

        stage('4. Functional Tests') {
            steps {

                sh '''
                    set -e

                    echo "===== CLEAN OLD CONTAINERS ====="

                    docker-compose down || true

                    docker rm -f shipping_mongo shipping_postgres || true
                '''

                sh '''
                    set -e

                    echo "===== START DATABASE ====="

                    docker-compose up -d postgres mongodb

                    echo "Waiting database startup..."
                    sleep 20

                    docker ps
                '''

                sh '''
                    set -e

                    echo "===== FUNCTIONAL TEST ====="

                    docker run --rm \
                        --network host \
                        -v $WORKSPACE:/app \
                        -w /app \
                        golang:1.25 \
                        sh -c "
                            set -e

                            if [ ! -f go.mod ]; then
                                echo 'ERROR: go.mod not found'
                                exit 1
                            fi

                            go mod tidy

                            go test -v -tags=functional ./internal/repository/...
                        "
                '''
            }

            post {
                always {
                    sh '''
                        docker-compose down || true
                    '''
                }
            }
        }

        stage('5. Push to Docker Hub') {
            steps {
                withCredentials([
                    usernamePassword(
                        credentialsId: 'docker-hub-id',
                        passwordVariable: 'DOCKER_PASS',
                        usernameVariable: 'DOCKER_USER_ENV'
                    )
                ]) {

                    sh '''
                        set -e

                        echo "$DOCKER_PASS" | docker login \
                            -u "$DOCKER_USER_ENV" \
                            --password-stdin

                        docker push ${DOCKER_IMAGE}
                    '''
                }
            }
        }

        stage('6. Deploy to K8s') {
            steps {
                withKubeConfig([credentialsId: 'kubeconfig-id']) {

                    sh '''
                        set -e

                        echo "===== DEPLOY TO K8S ====="

                        kubectl apply -f k8s/deployment.yaml

                        kubectl apply -f k8s/service.yaml

                        kubectl get pods
                    '''
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
        }

        success {
            echo 'Pipeline SUCCESS.'
        }

        failure {
            echo 'Pipeline FAILED.'
        }
    }
}