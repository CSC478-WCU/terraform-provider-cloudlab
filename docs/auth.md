# Authentication & Credentials (CloudLab)

> This page shows how to download your CloudLab credentials and convert the encrypted key into a decrypted PEM you can use for automation/SSH.

---

## 1) Download your credentials

1. Go to the CloudLab homepage and **sign in**.  
2. Click your **username** in the top-right → **Download Credentials**.  
   ![](assets/images/download_credentials.png)
3. On the credentials page, click **“here”** to download.  
   ![](assets/images/click_here.png)

You’ll get a file named **`cloudlab.pem`** that looks like this:

```
-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,DB055C9AE5A0F021

##############################################
##############################################
##############################################

-----END RSA PRIVATE KEY-----
-----BEGIN CERTIFICATE-----

##############################################
##############################################
##############################################

-----END CERTIFICATE-----
```

---

## 2) Decrypt the key (make a usable PEM)

Run this one-liner to replace the encrypted private-key block with the decrypted one and write to **`cloudlab_decrypted.pem`**:

```bash
perl -0777 -pe 'BEGIN{$dec=`openssl rsa -in cloudlab.pem`} s/-----BEGIN RSA PRIVATE KEY-----.*?-----END RSA PRIVATE KEY-----/$dec/s' cloudlab.pem > cloudlab_decrypted.pem
```

> You’ll be prompted for your cloudlab passphrase. The resulting `cloudlab_decrypted.pem` keeps the certificate block intact and swaps in a decrypted private key.

The decrypted file will look like:

```
-----BEGIN PRIVATE KEY-----
##############################################
##############################################
##############################################

-----END PRIVATE KEY-----

-----BEGIN CERTIFICATE-----
##############################################
##############################################
##############################################

-----END CERTIFICATE-----
```

Use **`cloudlab_decrypted.pem`** as your SSH key path in Terraform or Ansible.
