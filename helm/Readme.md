# How to deploy helm chart

`helm install ranutif helm/ --create-namespace --namespace ranutif `

then add api as a secret

`kubectl create secret generic tgapi --from-literal=key=<api> -n ranutif`

check it 

`kubectl get secret tgapi -n ranutif --template={{.data.key}} | base64 -D `

# For update with a tag

`helm upgrade ranutif helm/ --install --set image.tag=0.0.1`

