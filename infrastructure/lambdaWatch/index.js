const aws = require('aws-sdk');
const sqs = new aws.SQS({ apiVersion: '2012-11-05' });
const ecs = new aws.ECS({ apiVersion: '2014-11-13' });

const queue = process.env.SQS;
const taskDefinition = process.env.WORKERNAME;
const cluster = process.env.CLUSTER;
const SUBNET = process.env.SUBNET;

exports.handler = async (event, context) => {
  const p1 = new Promise((resolve, reject) => {
    const params = {
      MessageBody: JSON.stringify(event),
      QueueUrl: queue
    };
    sqs.sendMessage(params, (err, data) => {
      if (err) {
        console.warn('Error while sending message: ' + err);
        reject(err);
      }
      else {
        console.info('Message sent, ID: ' + data.MessageId);
        resolve(data);
      }
    });
  });
  try {
    await p1;
  } catch (err) {
    return context.fail(`SQS Fail: ${err}`);
  }

  const p2 = new Promise((resolve, reject) => {
    const params = {
      cluster,
      taskDefinition,
      launchType: 'FARGATE',
      count: 1,
      networkConfiguration: {
        awsvpcConfiguration: {
          assignPublicIp: "ENABLED",
          subnets: [SUBNET]
        }
      },
    };
    ecs.runTask(params, (err, data) => {
      if (err) {
        console.warn('error: ', "Error while starting task: " + err);
        reject(err);
      }
      else {
        console.info('Task started: ' + JSON.stringify(data.tasks))
        resolve(data);
      }
    });
  });

  try {
    await p2;
  } catch (err) {
    return context.fail(`ECS Fail: ${err}`);
  }
  context.succeed('Successfully processed Amazon S3 URL.');
};
