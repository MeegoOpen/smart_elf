const {
  execSync
} = require('child_process');
const fs = require('fs');
const inquirer = require('inquirer');

const distManifest = 'plugin.config.json';

const runCmd = (cmd) => {
  execSync(cmd, {
    stdio: 'inherit'
  }); // ignore_security_alert
};


(async () => {
  try {
    let [, , action, token] = process.argv;
    if (!action) {
      const answers = await inquirer.prompt([{
        type: 'list',
        name: 'action',
        message: 'Please select "action":',
        default: 2,
        choices: ['start', 'build', 'deploy'],
      }, ]);

      action = answers.action;
      console.log('Action is: ', action);
    }
    if (!fs.existsSync(distManifest)) {
      throw new Error(`Missing ${distManifest}. Please ensure it exists in project root.`);
    }
    const manifestConfig = JSON.parse(fs.readFileSync(distManifest, 'utf8'));
    const { siteDomain, pluginId } = manifestConfig || {};

    switch (action) {
      case 'deploy':
        runCmd(`lpm release ${token}`);
        console.log(
          '\nPlease goto here to deploy: \n\n\x1b[36m%s\x1b[0m\n',
          `${siteDomain}/openapp/${pluginId}#versions`,
        );
        break;
      default:
        runCmd(`lpm ${action}`);
    }
  } catch (err) {
    console.error(err);
    process.exit(1);
  }
})();
