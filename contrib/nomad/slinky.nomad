job "slinky-dev" {
  type        = "service"
  datacenters = ["skip-nomad-aws-us-east-2"]
  region      = "global"

  namespace = "slinky-dev"

  group "sidecar" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "slinky-sidecar-dev-http-service"
      port = "8080"

      connect {
        sidecar_service {}
      }

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-sidecar-dev-http-service.rule=Host(`slinky-sidecar-dev-http.skip.money`)",
        "traefik.http.routers.slinky-sidecar-dev-http-service.entrypoints=http",
      ]
    }

    service {
      name = "slinky-sidecar-dev-metrics-service"
      port = "8002"

      connect {
        sidecar_service {}
      }

      tags = [
        "metrics",
        "logs.promtail=true",
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-sidecar-dev-metrics-service.rule=Host(`slinky-sidecar-dev-metrics.skip.money`)",
        "traefik.http.routers.slinky-sidecar-dev-metrics-service.entrypoints=http",
      ]
    }

    service {
      name = "slinky-sidecar-dev-pprof-service"
      port = "6060"

      connect {
        sidecar_service {}
      }

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-sidecar-dev-pprof-service.rule=Host(`slinky-sidecar-dev-pprof.skip-internal.money`)",
        "traefik.http.routers.slinky-sidecar-dev-pprof-service.entrypoints=internal",
      ]
    }

    task "sidecar" {
      driver = "docker"

      config {
        image = "[[ .sidecar_image ]]"
        force_pull = true
        entrypoint = ["slinky", "--oracle-config-path", "/etc/slinky/default_config/oracle.json", "--market-config-path", "/etc/slinky/default_config/market.json"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }

  }

  group "chain" {
    count = 1

    network {
      mode = "bridge"
    }

    service {
      name = "slinky-simapp-dev-rpc-service"
      port = "26657"

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-simapp-dev-rpc-service.rule=Host(`slinky-simapp-dev-rpc.skip-internal.money`)",
        "traefik.http.routers.slinky-simapp-dev-rpc-service.entrypoints=internal",
      ]

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "slinky-sidecar-dev-http-service"
              local_bind_port  = 8080
            }
          }
        }
      }
    }

    service {
      name = "slinky-simapp-dev-lcd-service"
      port = "1317"

      connect {
          sidecar_service {}
      }

      tags = [
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-simapp-dev-lcd-service.rule=Host(`slinky-simapp-dev-lcd.skip-internal.money`)",
        "traefik.http.routers.slinky-simapp-dev-lcd-service.entrypoints=internal",
      ]
    }

    service {
      name = "slinky-simapp-dev-chain-metrics-service"
      port = "26660"

      connect {
          sidecar_service {}
      }

      tags = [
        "metrics",
        "logs.promtail=true",
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-simapp-dev-chain-metrics-service.rule=Host(`slinky-simapp-dev-chain-metrics.skip.money`)",
        "traefik.http.routers.slinky-simapp-dev-chain-metrics-service.entrypoints=http",
      ]
    }

    service {
      name = "slinky-simapp-dev-app-metrics-service"
      port = "8001"

      connect {
          sidecar_service {}
      }

      tags = [
        "metrics",
        "logs.promtail=true",
        "traefik.enable=true",
        "traefik.consulcatalog.connect=true",
        "traefik.http.routers.slinky-simapp-dev-app-metrics-service.rule=Host(`slinky-simapp-dev-app-metrics.skip.money`)",
        "traefik.http.routers.slinky-simapp-dev-app-metrics-service.entrypoints=http",
      ]
    }

    volume "data" {
      type            = "csi"
      read_only       = false
      source          = "slinky-simapp-dev-node-volume"
      access_mode     = "single-node-writer"
      attachment_mode = "file-system"
    }

    task "init" {
      driver = "docker"

      volume_mount {
        volume      = "data"
        destination = "/src/slinky/tests/.slinkyd"
        read_only   = false
      }

      config {
        image      = "[[ .chain_image ]]"
        force_pull = true
        entrypoint = ["sh", "-c", "/tmp/init.sh"]
        volumes    = ["local/tmp/init.sh:/tmp/init.sh"]
      }

      template {
        data = <<EOH
#!/bin/sh
rm -rf tests/.slinkyd/**

make build-configs
sed -i 's\oracle:8080\localhost:8080\g' /src/slinky/tests/.slinkyd/config/app.toml
        EOH

        perms = "777"

        destination = "local/tmp/init.sh"
      }

      lifecycle {
        hook    = "prestart"
        sidecar = false
      }
    }

    task "chain" {
      driver = "docker"

      volume_mount {
        volume      = "data"
        destination = "/src/slinky/tests/.slinkyd"
        read_only   = false
      }

      config {
        image      = "[[ .chain_image ]]"
        force_pull = true
        entrypoint = ["make", "start-app"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
