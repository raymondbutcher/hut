# Hut

Status: pre-alpha

Hut is a very simple Terraform runner that works with a particular standard Terraform project structure to automatically add appropriate `-chdir`, `-var-file`, and `-backend-config` arguments before running Terraform.

* Hut passes all command line arguments through to Terraform, so commands like `hut plan -target=ADDRESS` work just fine. If the Hut project structure is not detected, it's just like using Terraform directly.
* Hut does not create any temporary directories or temporary files; all it does is run Terraform.
* Hut prints out the full Terraform command before running it, so you can always see what Hut is doing.

## Project structure

If Hut detects this structure, then it adds extra arguments before running Terraform. If this project structure is not detected, then no arguments get added and it's just like running Terraform directly.

This structure is [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself). It allows a [root module](https://www.terraform.io/language/modules#the-root-module) to be deployed to multiple environments by putting environment-specific variables and backend configuration in environment-specific subdirectories. Hut figures out the extra Terraform arguments needed to make Terraform
work with this directory structure.

This directory structure is still valid/standard/vanilla Terraform. Terraform can still be used directly if passed in the same arguments would be added by Hut. Hut even prints out the final Terraform command before it runs, so it can be copied and used directly (e.g. in CI/CD scripts).

Specification:

* The [root module](https://www.terraform.io/language/modules#the-root-module) must contains at least one `*.tf[.json]` file.
* The root module may contain any number of subdirectories and they may be nested (multiple levels deep).
* The root module, or any of its subdirectories, may contain `terraform.tfvars[.json]`, `*.auto.tfvars[.json]`, and `*.tfbackend` files. There may be any number of these, and in any location.
* When Hut runs from a root module directory, or any of its subdirectories, and the directory contains a `terraform.tfvars[.json]` file:
  * If running from a subdirectory, `TERRAFORM_DATA_DIR=${current_dir}/.terraform` and `-chdir=${root_module_dir}` are added to the Terraform command so it runs in the root module directory while still using the current directory for the `.terraform` data directory.
  * For `terraform.tfvars[.json]`, `*.auto.tfvars[.json]`, and `*.tfbackend` files in the current directory and any parent directories, stopping at the root module, a `-var-file` or `-backend-config` argument is added to the Terraform command as appropriate.

### What is it like to use?

Have you noticed how the Terraform commands [in](https://developer.hashicorp.com/terraform/tutorials/cli/init) [most](https://developer.hashicorp.com/terraform/tutorials/cli/plan) [tutorials](https://developer.hashicorp.com/terraform/tutorials/configuration-language/variables) don't include any `-var-file` or `-backend-config` arguments? You write some code and run short commands like `terraform init` and `terraform plan`. Terraform seems so nice and simple to use!

But things are complicated when you introduce multiple environments. There are a few options:

* Maintain multiple root modules, one for each environment.
  * Introduces code duplication which can become painful as the project grows.
* Use [workspaces](https://www.terraform.io/language/state/workspaces).
  * Requires that all environments share the same Terraform backend.
  * The directory becomes stateful with the currently selected workspace.
* Manually add `-var-file`, `-backend-config`, `-chdir` arguments when running Terraform.
  * Easy to forget.
  * Easy to make a mistake.
  * More stressful if you're not confident that you have all of the correct arguments.
* Use a Terraform wrapper.
  * There are many options out there.
  * Do they require much from you? (learning, installation, dependencies, configuration, etc)
  * Do they lock you in and become difficult to stop using?

Using Hut is like those simple Terraform tutorials, except that you can have multiple environments just by putting environment files in environment directories and running Hut from those directories.

Day-to-day usage might look like this:

```
~ # cd example
~/example # hut fmt -recursive

~/example # cd au/dev
~/example/au/dev # hut init
~/example/au/dev # hut plan
~/example/au/dev # hut apply -target=module.example

~/example/au/dev # cd ../prod
~/example/au/prod # hut init
~/example/au/prod # hut plan

~/example/au/prod # cd ../../eu/dev
~/example/eu/dev # hut init
~/example/eu/dev # hut plan

~/example/eu/dev # cd ../../eu/prod
~/example/eu/prod # hut init
~/example/eu/prod # hut plan
```

Notice how the commands are short, and the current directory makes the target environment very obvious.

## Lineage

The author of Hut has built a number of Terraform wrappers over the years, and Hut is just the latest one.

* 2018: [Jinjaform](https://github.com/claranet/jinjaform)
  * Status: discontinued
  * Language: Python
  * Motives:
    * Allow DRY project structure.
    * Use Jinja2 templates in Terraform projects.
    * Add support for AWS MFA prompts (Terraform did not at the time).
  * Learned:
    * Building a temporary directory and running Terraform from there causes problems with relative paths, and makes it hard to run Terraform directly.
    * Jinja2 and Terraform syntax do not mix well.
* 2019: [Pretf](https://github.com/raymondbutcher/pretf)
  * Status: stable
  * Language: Python
  * Motives:
    * Allow more flexible project structures and customisation.
    * Don't build/run Terraform from a temporary directory.
    * Generate Terraform code with Python or Jinja2.
    * Add support for AWS MFA prompts (Terraform did not at the time).
  * Learned:
    * People can see flexibility as it being complicated and requiring learning/effort to start using, even if the provided examples do everything they need.
    * Python is not for everyone. Pretf only requires a single small Python file to enable and configure it, but it still turns your Terraform project into a Python project.
    * Python dependencies are almost always undesireable, e.g. needing a Python + Terraform Docker image in CI/CD jobs, and needing to use Pip/Poetry/etc. to install Pretf.
* 2022: [LTF](https://github.com/raymondbutcher/ltf)
  * Status: alpha
  * Language: Go
  * Motives:
    * Release as single binary.
    * Keep it simple and a little flexible.
    * Mostly compatible with vanilla/standard Terraform.
  * Learned:
    * When it's nearly compatible with vanilla/standard Terraform then it might be worth going the extra mile to make it 100% compatible, even if you lose some features.
* 2022: Hut (this project)
  * Status: pre-alpha
  * Language: Go
  * Motives:
    * Make it super simple.
    * 100% compatible with vanilla/standard Terraform.
  * Learned:
    * Unsure if it's worth having 2 very similar projects (LTF and Hut). Maybe I should just remove some features from LTF, or make LTF automatically run in a more simple mode when certain features (hooks) aren't being used.
