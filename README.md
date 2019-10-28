# FillPDF

FillPDF is a golang library to easily fill PDF forms. This library uses the pdftk utility to fill the PDF forms with fdf data.
Currently this library only supports PDF text and checkbox field values. Feel free to add support to more form types (Send pull request to original developer)
This fork extends with some more pdftk commands

- Multistamp
- Ability to generate PDF's with special characters (with flatten) with pdftk. (Limited by font in PDF)

2019: the library has been turned into a standalone server.

## Build and push container

After having merged your changes to the master branch of this repo, build and push a new container:

`$ gcloud builds submit -t gcr.io/$PROJECT/fillpdf:v$VERSION --project $PROJECT`

Change `$PROJECT` and `$VERSION` to what you need. For the version numbers, use semver principles.

NOTE that the container must be built and pushed for all the projects where you want to deploy the service. Attow these are:

- `mindoktor-dev` (dev and test)
- `e-vard` (sweden prod)
- `docly-prod` (UK prod)

## Redeploy the service

During the deployment steps, it is _vital_ that you have selected the correct context and namespace using `kubectx` and `kubens`. Otherwise you may deploy to the wrong project/environment. When deploying to test, note that there is a single fillpdf service for all dev/test envs. This is found in the `default` namespace of the `test-mindoktor-dev` context.

### Step 1

In the `mindoktor` repo root: Bump the version in the fillpdf deployment file: `devops/k8s/deployment_fillpdf.yaml`

### Step 2

Then generate env files and config for the SITE you wish to deploy to. For example:

`$ SITE=se-test fab gen_env_files`

After the above, a large number of files will have been modified in your branch. This does not matter as you will only be deploying the changes to one of these.

### Step 3

Before this, double check that you have the correct context and namespace selected. Also, if you want to monitor the progress of the deployment, i.e. when the new pod goes up and the old one down, you can monitor the state of the pod in a separate terminal like so: `watch 'kubectl get pods| grep fillpdf'`

Then:

`$ kubectl apply -f devops/k8s/deployment_fillpdf.yaml`

This should make a new fillpdf pod deploy and the old one be terminated.

## Documentation

General docs for the fillpdf library: [GoDoc.org](https://godoc.org/github.com/desertbit/fillpdf).
