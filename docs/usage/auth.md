# Authentication with PEM

CloudLab uses **mutual TLS** to its XML-RPC portal. You’ll download credentials and (optionally) decrypt the private key for automation.

!!! warning 

    **Please keep the credentials a secret** in Github Actions or repository workflows — anyone with the decrypted .pem can access cloudlab's infrastructure and every user in your project's active experiments. 

---

!!! example "Step 1 — Download credentials"

    1. Sign in to CloudLab  
    2. Click your **username → Download Credentials**  
    3. Click **“here”** to download

    ![Download credentials](../assets/images/download_credentials.png)

    ![Click here](../assets/images/click_here.png)

---

!!! example "Step 2 — Decrypt private key"

    ```bash
    perl -0777 -pe 'BEGIN{$dec=`openssl rsa -in cloudlab.pem`} s/-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----/$dec/s' cloudlab.pem > cloudlab_decrypted.pem
    ```

    You’ll be prompted for your passphrase. The output will contain:

    ```
    -----BEGIN PRIVATE KEY-----
    ...
    -----END PRIVATE KEY-----

    -----BEGIN CERTIFICATE-----
    ...
    -----END CERTIFICATE-----
    ```

Use ``cloudlab_decrypted.pem`` as the PEM path in Terraform or Ansible. 

