const goLibrary = require('go-library');

const options = {
  destinationPath: 'bin',
  repo: 'jurekbarth/pup',
  version: 'v0.0.1',
  projectname: 'pup'
}


goLibrary(options);
