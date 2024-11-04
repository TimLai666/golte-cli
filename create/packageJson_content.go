package create

const packageJsonContentTemplate = `{
  "name": "{{projectName}}",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "type": "module",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "keywords": [],
  "author": "",
  "license": "ISC",
  "dependencies": {
    "golte": "^0.1.1"
  }
}
`
