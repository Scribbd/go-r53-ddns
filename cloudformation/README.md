# How to use the Cloudformation template

1. Deploy this stack within Cloudformation
1. Go to the IAM dashboard and generate the accesskeys
1. Optional if you are not me: delete the vault.yml file and create a new one with `ansible-vault create vault.yml`
1. Store the accesskeys in the Ansible Vault vault.yml file

If you are going to use the ansible playbook to rotate the secrets. You won't need to generate the access keys for the R53-DDNS k8s deployment.