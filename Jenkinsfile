// Pipelines
pipeline {
  agent any
  environment {
    dockerHub = "harbor.clevabit.services"
    imageName = "ci/flagr"
   }
  stages {
    stage('Build') {
      steps {
        sh '''
            docker build -t ${dockerHub}/${imageName}:${GIT_BRANCH#remotes/origin/}-${BUILD_NUMBER} .
            docker push ${dockerHub}/${imageName}:${GIT_BRANCH#remotes/origin/}-${BUILD_NUMBER}
        '''
      }
    }
    stage('Deploy to kubernetes to prod app for clevabit branch') {
      when {
        allOf {
          branch 'clevabit'
          expression {
            currentBuild.result == null || currentBuild.result == 'SUCCESS'
          }
        }
      }
      steps {
        kubernetesDeploy(kubeconfigId: 'prod-app-we-kubeconfig',
            configs: 'deployment.yaml')
      }
    }
  }
}
