prom2owl
=========

A tool to convert prometheus data format to owl.

# Usage

Installing and building:

    $ go get github.com/owl-plugins/prom2owl

Running:

    $ prom2owl http://my-prometheus-client.example.org:8080/metrics

Running with TLS client authentication:

    $ prom2owl -cert=/path/to/certificate -key=/path/to/key http://my-prometheus-client.example.org:8080/metrics

Running with custom labels:
    # prom2owl -labels tagk1=tagv1,tagk2=tagv2 http://my-prometheus-client.example.org:8080/metrics
