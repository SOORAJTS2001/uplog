import { useEffect } from "react";

function CLIDownloadRedirect() {

  useEffect(() => {
    const downloadUrl = `https://raw.githubusercontent.com/SOORAJTS2001/uplog/main/install.sh`;
    window.location.href = downloadUrl;
  },[]);

  return <p>Redirecting to cli download…</p>;
}

export default CLIDownloadRedirect;
