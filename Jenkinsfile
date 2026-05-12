pipeline {
    agent any

    environment {
        IMAGE_NAME = 'dickaf/shipping-service'
        IMAGE_TAG  = "${BUILD_ID}"
    }

    stage('Docker Test') {
    steps {
        bat 'docker pull golang:1.25-alpine'
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