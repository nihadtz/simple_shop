# Welcome to SimpleShop!

This application provides simple shop using JWT for authentication, Stripe for payments, RBAC for access control, and allows database migrations.

# How it works

Every user is able to visit public dashboard and see the most popular and all products.

**Endpoint** 
- /dashboard/public 

If user is interested in a purchase of a products first login is a must have.

**Endpoint** 
- /register

Next step is to login and checking if everything is ok with profile.

**Endpoints** 
 - /login
 - /user/:id - only his account

After login user is able to list all products and purchase one of them.

**Endpoints** 
 - /products
 - /product/:id
 - /purchase

After done some purchases user is able to list all his purchases.

**Endpoint** 
 - /purchases
 - /purchase/:id - only his account

To make purchase completed user can pay it using Stripe.

**Endpoint** 
 - /payment/stripe/:id  - purchase ID

Other type of user is a admin, which has access to everything mention above, except that he can see all users, purchases and additionally admin dashboard. Only thing which admin can't do is payment for other user purchase.

**Endpoint** 
 - /dashboard/admin

Admin dashboard is providing insight about product sales, in following manner:

**List of all products with following attributes**
 - SoldQuantity is total quantity sold in all Purchases
 - Income is total value of all Purchases
 - Paid is amount successfully paid by Customers
 - NotPaid is difference between Income and Paid
 - FailedPayments is amount failed during payment process

Admin is also able to add new products and to edit them.

**Endpoint** 
 - /product