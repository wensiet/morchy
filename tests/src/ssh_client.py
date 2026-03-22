from dataclasses import dataclass
import paramiko


@dataclass
class SSHResult:
    exit_code: int
    stdout: str
    stderr: str


class SSHClient:
    def __init__(
        self,
        host: str,
        username: str,
        password: str | None = None,
        key_filename: str | None = None,
        port: int = 22,
        timeout: int = 10,
    ):
        self.host = host
        self.username = username
        self.password = password
        self.key_filename = key_filename
        self.port = port
        self.timeout = timeout

    def run(self, command: str, get_pty: bool = False) -> SSHResult:
        client = paramiko.SSHClient()
        client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

        try:
            client.connect(
                hostname=self.host,
                port=self.port,
                username=self.username,
                password=self.password,
                key_filename=self.key_filename,
                timeout=self.timeout,
            )

            stdin, stdout, stderr = client.exec_command(command, get_pty=get_pty)

            exit_code = stdout.channel.recv_exit_status()
            out = stdout.read().decode()
            err = stderr.read().decode()

            return SSHResult(exit_code=exit_code, stdout=out, stderr=err)

        finally:
            client.close()
