using System;
using System.Collections.Generic;
using System.Linq;
using System.Web;

namespace WebApplication1
{
    /// <summary>
    /// Handler1 的摘要说明
    /// </summary>
    public class Handler1 : IHttpHandler
    {

public void ProcessRequest(HttpContext context)
{
    if (context.Request.Files.Count > 0)
    {
        context.Request.Files[0].SaveAs("d:\\1.wav");
    }
}

        public bool IsReusable
        {
            get
            {
                return false;
            }
        }
    }
}