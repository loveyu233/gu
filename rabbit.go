package gu

import (
	"crypto/tls"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type RabbitMQConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Vhost    string
}

var DefaultRabbitMQConfig = RabbitMQConfig{
	Host:     "localhost",
	Port:     5672,
	Username: "rabbit",
	Password: "rabbit123",
}

func MustInitDefaultRabbitMQ() *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", DefaultRabbitMQConfig.Username, DefaultRabbitMQConfig.Password, DefaultRabbitMQConfig.Host, DefaultRabbitMQConfig.Port)

	dial, err := amqp.Dial(url)
	if err != nil {
		logrus.Errorf("Failed to connect to RabbitMQ at host %s. Error: %s", url, err)
		return nil
	}
	RabbitClient = dial
	return RabbitClient
}

func MustInitRabbitMQ(config RabbitMQConfig, tlsConfig ...*tls.Config) (*amqp.Connection, error) {
	var scheme = "amqp"
	if len(tlsConfig) > 0 {
		scheme = "amqps"
	}

	url := amqp.URI{
		Scheme:   scheme,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Vhost:    config.Vhost,
	}.String()

	var err error

	if scheme == "amqps" {
		RabbitClient, err = amqp.DialTLS(url, tlsConfig[0])
		if err != nil {
			logrus.Errorf("Failed to connect to RabbitMQ at host %s. Error: %s", url, err)
			return nil, err
		}
	} else {
		RabbitClient, err = amqp.Dial(url)
		if err != nil {
			logrus.Errorf("Failed to connect to RabbitMQ at host %s. Error: %s", url, err)
			return nil, err
		}
	}

	return RabbitClient, nil
}

func RabbitMQGetMsg(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	declare, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	consume, err := channel.Consume(declare.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return consume, nil
}

func RabbitMQSendMsg(channel *amqp.Channel, queueName string, body []byte) error {
	declare, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{Body: body}
	err = channel.Publish("", declare.Name, false, false, msg)
	if err != nil {
		return err
	}

	return nil
}
