# FillPDF

FillPDF is a golang library to easily fill PDF forms. This library uses the pdftk utility to fill the PDF forms with fdf data.
Currently this library only supports PDF text and checkbox field values. Feel free to add support to more form types (Send pull request to original developer)
This fork extends with some more pdftk commands

- Multistamp
- Ability to generate PDF's with special characters (with flatten) with pdftk. (Limited by font in PDF)

2019: the library has been turned into a standalone server.

Build and push a new container like this:
`$ gcloud builds submit -t gcr.io/$PROJECT/fillpdf:v$VERSION --project $PROJECT`
Change $PROJECT and $VERSION to what you need.

Then you have to redeploy the service. Bump the version in main repo: devops/k8s/deployment_fillpdf.yaml
For example:
$ SITE=se-test fab gen_env_files
$ kubectl apply -f devops/k8s/deployment_fillpdf.yaml

## Documentation

Check the Documentation at [GoDoc.org](https://godoc.org/github.com/desertbit/fillpdf).
