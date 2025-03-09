# cribl
terraform provider for cribl.io


Add Pipeline

```mermaid
sequenceDiagram
participant pipeline
participant functions
participant inputs
participant outputs
participant routes
participant pipelines
pipeline --> functions: fetch
pipeline --> inputs: fetch
pipeline --> outputs: fetch
pipeline --> routes: fetch
pipeline -> pipelines: create pipeline<br>(basically a stub)
```
