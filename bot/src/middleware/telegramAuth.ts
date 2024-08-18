import { NextFunction, Request, Response } from "express";
import { validateFromTelegram } from "../utils/validateFromTelegram";

export const telegramAuth = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  if (!req.headers.authorization) return res.status(403).end();
  try {
    const isValid = validateFromTelegram(req.headers.authorization || "");
    if (!isValid) return res.status(403).end();
    res.locals.user = isValid;
    next();
  } catch (error) {
    res.status(403).end();
  }
};
