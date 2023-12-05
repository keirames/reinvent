import express, { Request, Response } from 'express';

const app = express();

export const main = () => {
  app.get('/', (req, res) => {
    console.log(`from instance No.${process.env.NO}`);
    res.send('hello');
  });

  app.listen(process.env.PORT, () => {
    console.log(`listen on port ${process.env.PORT}`);
  });
};

main();
