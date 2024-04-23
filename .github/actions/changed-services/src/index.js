const path = require("path");
const core = require("@actions/core");
const { readdir } = require("fs/promises");
const { exec } = require("child_process");

const baseRepository = core.getInput("base-repository", {
  required: true,
  trimWhitespace: true,
});

const getDirectories = async (source) =>
  (await readdir(source, { withFileTypes: true }))
    .filter((dirent) => dirent.isDirectory())
    .map((dirent) => dirent.name);

const getDependencyGraph = async (service) => {
  return new Promise((resolve, reject) => {
    exec(`go list -json cmd/${service}/main.go`, (error, stdout, stderr) => {
      if (error) {
        return reject(error);
      }
      if (stderr) {
        return reject(stderr);
      }
      return resolve({
        service,
        dependencies: JSON.parse(stdout)
          ["Deps"].filter((dep) => dep.includes(baseRepository))
          .map((dep) => dep.replace(`github.com/${baseRepository}/`, "")),
      });
    });
  });
};

async function main() {
  try {
    const changedFiles = JSON.parse(
      core.getInput("changed-files", { required: true, trimWhitespace: true }),
    );
    const services = await getDirectories(
      path.join(__dirname, "../../../../services"),
    );

    core.debug(`Changed files: ${changedFiles.join(", ")}`);
    if (changedFiles.includes("go.mod") || changedFiles.includes("go.sum")) {
      core.setOutput("services_count", services.length);
      core.debug(`Changed services: ${services.join(", ")}`);
      return core.setOutput("services", JSON.stringify(services));
    }

    const dependencies = (
      await Promise.all(services.map(getDependencyGraph))
    ).reduce(
      (acc, { service, dependencies }) => ({ ...acc, [service]: dependencies }),
      {},
    );

    const needsBuild = [];
    for (let index = 0; index < services.length; index++) {
      const service = services[index];
      const serviceDependencies = [
        `services/${service}`,
        ...dependencies[service],
      ];
      for (let fileIndex = 0; fileIndex < changedFiles.length; fileIndex++) {
        const file = changedFiles[fileIndex];
        if (serviceDependencies.some((dep) => file.includes(dep))) {
          needsBuild.push(service);
          break;
        }
      }
    }

    core.debug(`Changed services: ${needsBuild.join(", ")}`);
    core.setOutput("services_count", needsBuild.length);
    return core.setOutput("services", JSON.stringify(needsBuild));
  } catch (error) {
    core.debug(error);
    core.setFailed(error.message);
  }
}

main();
